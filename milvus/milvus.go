package milvus

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/types"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	// 主键
	FiledPK = "pk"
	// 切分块 ID
	FiledID = "fid"
	// 块内容
	//FiledSource = "source"
	// 块向量
	FiledEmbedding = "embedding"
	// 文档 ID
	FiledDocumentID = "doc"
	// 集合名称
	FiledCorpus = "corpus"
)

type Config struct {
	// Addr is the address of the Milvus server.
	Addr string

	// Username
	User string

	// Password
	Password string
}

type CollectionConfig struct {

	// database , 最大 64 个, 现在未启用
	Dadabase string `json:"database"`

	// partition name， 分区，现在未启用
	// Milvus 允许将大量矢量数据划分为少量分区。然后，搜索和其他操作可以限制在一个分区，以提高性能。
	PartitionName string `json:"partitionName"`

	// collection belong to Corpus
	Corpus string `json:"corpus"`

	// collection name
	CollectionName string `json:"collectionName"`

	// collection description,  默认是 `Collection for ${Corpus}`
	CollectionDescription string `json:"collectionDescription"`

	// index and search cluster units for collection， default is 3
	Units int `json:"units" default:"10"`

	// 分区
	Shards int32 `json:"shards" default:"0"`

	// embedding model
	Model string `json:"model" default:"m3e"`

	// metric type, 默认欧几里得距离
	MetricType string `json:"metricType" default:"IP"`

	// Index type, default is IVF_FLAT by faiss， 目前只支持 IVF_FLAT
	IndexType string `json:"indexType" default:"IVF_FLAT"`

	// 默认 768
	Dimension int `json:"dimension" default:"768"`
}

/*
## Limits
|Feature|Maximum limit|
|---|---|
|Length of a collection name|255 characters|
|Number of partitions in a collection|4,096|
|Number of fields in a collection|256|
|Number of shards in a collection|256|
|Dimensions of a vector|32,768|
|Top K|16,384|
|Target input vectors|16,384|
*/

func (cfg *Config) init() {
	if cfg.Addr == "" {
		cfg.Addr = "localhost:19530"
	}
}

type Milvus struct {
	client client.Client
	ctx    context.Context
}

func DefalutConfig() CollectionConfig {
	p := CollectionConfig{}
	common.DefaultConfig(&p, "default")
	return p
}

func ParseConfig(conf map[string]interface{}) (CollectionConfig, error) {
	p := DefalutConfig()
	j, err := json.Marshal(conf)
	err = json.Unmarshal(j, &p)
	return p, err
}

func NewMilvus(cfg *Config) (*Milvus, error) {
	cfg.init()
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	//defer cancel()
	c, err := client.NewDefaultGrpcClientWithURI(ctx, cfg.Addr, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}

	m := &Milvus{
		client: c,
		ctx:    ctx,
	}
	return m, nil
}

func (m *Milvus) Insert(collection CollectionConfig, chunks map[string][]*types.Chunk) error {
	var idList []string
	//var textList []string
	var documentIDList []string
	var embeddingList [][]float32
	var corpusIDList []string

	for _, chunkList := range chunks {
		for _, chunk := range chunkList {
			idList = append(idList, chunk.ID)
			//textList = append(textList, chunk.Text)
			documentIDList = append(documentIDList, chunk.DocumentID)
			corpusIDList = append(corpusIDList, chunk.Metadata.Corpus)
			embeddingList = append(embeddingList, chunk.Embedding)
		}
	}
	// make list
	idCol := entity.NewColumnVarChar(FiledID, idList)
	//textCol := entity.NewColumnVarChar(FiledSource, textList)
	documentIDCol := entity.NewColumnVarChar(FiledDocumentID, documentIDList)
	corpusIDCol := entity.NewColumnVarChar(FiledCorpus, corpusIDList)
	embeddingCol := entity.NewColumnFloatVector(FiledEmbedding, collection.Dimension, embeddingList)
	_, err := m.client.Insert(m.ctx, collection.CollectionName, collection.PartitionName, idCol, documentIDCol, corpusIDCol, embeddingCol)

	return err
}

// TODO: Modify document sugment
func (m *Milvus) Modify(collection CollectionConfig, chunk *types.Chunk) error {

	expr := fmt.Sprintf(`%s == "%s" && %s == "%s"`, FiledDocumentID, chunk.ID, FiledCorpus, collection.Corpus)
	result, err := m.client.Query(m.ctx, collection.CollectionName, nil, expr, []string{FiledPK})
	//m.client.ReleaseCollection(m.ctx, collection.CollectionName)

	if err != nil {
		fmt.Println(err)
		return err
	}
	var pkCol *entity.ColumnInt64
	for _, field := range result {
		if field.Name() == FiledPK {
			if c, ok := field.(*entity.ColumnInt64); ok {
				pkCol = c
			}
		}
	}
	if len(pkCol.Data()) == 0 {
		return nil
	}
	err = m.client.DeleteByPks(m.ctx, collection.CollectionName, collection.PartitionName, pkCol)
	return err
}

func (m *Milvus) Flush(collection CollectionConfig) error {
	e := m.client.Flush(m.ctx, collection.CollectionName, false)
	if e != nil {

	}
	return e
}

// Delete deletes the chunks belonging to the given documentIDs.
// As a special case, empty documentIDs means deleting all chunks.
func (m *Milvus) Delete(collection CollectionConfig, documentIDs ...string) error {
	// To delete all chunks, we drop the old collection and create a new one.
	if len(documentIDs) == 0 {
		if err := m.client.ReleaseCollection(m.ctx, collection.CollectionName); err != nil {
			return err
		}
		if err := m.client.DropCollection(m.ctx, collection.CollectionName); err != nil {
			return err
		}
		return nil
	}

	// 表达式
	expr := fmt.Sprintf(`%s in ["%s"] && %s == "%s"`, FiledDocumentID, strings.Join(documentIDs, `", "`), FiledCorpus, collection.Corpus)
	// 查询
	m.client.LoadCollection(m.ctx, collection.CollectionName, false)
	result, err := m.client.Query(m.ctx, collection.CollectionName, nil, expr, []string{FiledPK})
	//m.client.ReleaseCollection(m.ctx, collection.CollectionName)

	if err != nil {
		fmt.Println(err)
		return err
	}

	var pkCol *entity.ColumnInt64
	for _, field := range result {
		if field.Name() == FiledPK {
			if c, ok := field.(*entity.ColumnInt64); ok {
				pkCol = c
			}
		}
	}
	if len(pkCol.Data()) == 0 {
		return nil
	}
	err = m.client.DeleteByPks(m.ctx, collection.CollectionName, collection.PartitionName, pkCol)
	m.Flush(collection)
	return err
}

func (m *Milvus) Get(collection CollectionConfig, doc string) (*types.Similarity, error) {
	m.client.LoadCollection(m.ctx, collection.CollectionName, false)
	expr := ""
	var partition []string = nil
	if collection.PartitionName != "" {
		partition = []string{collection.PartitionName}
	}
	if collection.Corpus != "" {
		expr = fmt.Sprintf(`%s == "%s" and %s == "%s"`, FiledCorpus, collection.Corpus, FiledID, doc)
	}

	result, err := m.client.Query(
		m.ctx,                     // ctx
		collection.CollectionName, // CollectionName
		partition,                 // partitionNames
		expr,                      // expr
		[]string{FiledID, FiledDocumentID, FiledCorpus}, // outputFields
	)
	/*
		m.client.ReleaseCollection(m.ctx, collection.CollectionName)
		if err != nil {
			return nil, err
		}
	*/
	if len(result) == 0 {
		return nil, nil
	}

	res, err := constructResult(result)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res[0], nil
}

// Query searches similarities of the given embedding with default consistency level.
func (m *Milvus) Query(collection CollectionConfig, embedding types.Embedding, limit int) ([]*types.Similarity, error) {
	expr := ""

	if collection.Corpus != "" {
		expr = fmt.Sprintf(`%s == "%s"`, FiledCorpus, collection.Corpus)
	}

	vec2search := []entity.Vector{
		entity.FloatVector(embedding),
	}

	// 这一块需要和 index type 一致
	// 适用于完美准确性和相对较小数据集的 entity.NewIndexFlatSearchParam()
	// 设置集群查询数量, 目前只支持 IVF_FLAT
	sp, _ := entity.NewIndexIvfFlatSearchParam(collection.Units)

	m.client.LoadCollection(m.ctx, collection.CollectionName, false)

	var partition []string = nil
	if collection.PartitionName != "" {
		partition = []string{collection.PartitionName}
	}

	metric := entity.IP
	if collection.MetricType == "L2" {
		metric = entity.L2
	}

	if limit == 0 {
		limit = 3
	}
	opt := client.WithSearchQueryConsistencyLevel(entity.ClStrong)
	result, err := m.client.Search(
		m.ctx,                     // ctx
		collection.CollectionName, // CollectionName
		partition,                 // partitionNames
		expr,                      // expr
		[]string{FiledID, FiledDocumentID, FiledCorpus}, // outputFields
		vec2search,     // vectors
		FiledEmbedding, // vectorField
		metric,         // metricType
		limit,          // topK
		sp,             // sp
		opt,
	)
	/*
		m.client.ReleaseCollection(
			m.ctx,                     // ctx
			collection.CollectionName, // CollectionName
		)
	*/
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return constructSimilaritiesFromResult(&result[0])
}

func (m *Milvus) CreateCollection(collection CollectionConfig, createNew bool) error {
	has, err := m.client.HasCollection(m.ctx, collection.CollectionName)
	if err != nil {
		return err
	}

	if has && !createNew {
		return nil
	}

	if has {
		_ = m.client.DropCollection(m.ctx, collection.CollectionName)
	}

	if collection.CollectionDescription == "" {
		collection.CollectionDescription = "Collection for organize " + collection.Corpus
	}

	// The collection does not exist, so we need to create one.
	schema := GetSchema(collection)

	// collection 和 index 参数
	// Create collection with consistency level, which serves as the default search/query consistency level.
	// WithConsistencyLevel, WithPartitionNum,  WithCollectionProperty
	// https://milvus.io/docs/consistency.md
	if err := m.client.CreateCollection(m.ctx, schema, collection.Shards, client.WithConsistencyLevel(entity.ClBounded)); err != nil {
		return err
	}

	// Create index "IVF_FLAT".
	// Number of cluster units, Seach 的 nprobe 要等于或者小于此值
	metric := entity.IP
	if collection.MetricType == "L2" {
		metric = entity.L2
	}
	idx, err := entity.NewIndexIvfFlat(metric, collection.Units)
	if err != nil {
		return err
	}

	return m.client.CreateIndex(m.ctx, collection.CollectionName, FiledEmbedding, idx, false)
}

func (m *Milvus) DeleteCollection(collection CollectionConfig) error {

	return m.client.DropCollection(m.ctx, collection.CollectionName)

}

func (m *Milvus) Close() error {
	return m.client.Close()
}

func GetSchema(collection CollectionConfig) *entity.Schema {

	schema := &entity.Schema{
		CollectionName:     collection.CollectionName,
		Description:        collection.CollectionDescription,
		AutoID:             true,
		EnableDynamicField: true,
		Fields: []*entity.Field{
			{
				Name:       FiledPK,
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     FiledID,
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: fmt.Sprintf("%d", 65535),
				},
			},
			{
				Name:     FiledDocumentID,
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: fmt.Sprintf("%d", 65535),
				},
			},
			{
				Name:     FiledCorpus,
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: fmt.Sprintf("%d", 65535),
				},
			},
			{
				Name:     FiledEmbedding,
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					entity.TypeParamDim: fmt.Sprintf("%d", collection.Dimension),
				},
			},
		},
	}
	return schema
}

func constructResult(result client.ResultSet) ([]*types.Similarity, error) {
	var idCol *entity.ColumnVarChar
	var documentIDCol *entity.ColumnVarChar
	var corpusIDCol *entity.ColumnVarChar
	for _, field := range result {
		switch field.Name() {
		case FiledID:
			if c, ok := field.(*entity.ColumnVarChar); ok {
				idCol = c
			}
		case FiledDocumentID:
			if c, ok := field.(*entity.ColumnVarChar); ok {
				documentIDCol = c
			}
		case FiledCorpus:
			if c, ok := field.(*entity.ColumnVarChar); ok {
				corpusIDCol = c
			}
		}
	}
	var similarities []*types.Similarity
	for i := 0; i < idCol.Len(); i++ {
		id, err := idCol.GetAsString(i)
		if err != nil {
			return nil, err
		}
		documentID, err := documentIDCol.GetAsString(i)
		if err != nil {
			return nil, err
		}
		corpusID, err := corpusIDCol.GetAsString(i)
		if err != nil {
			return nil, err
		}

		similarities = append(similarities, &types.Similarity{
			Chunk: &types.Chunk{
				ID:         id,
				DocumentID: documentID,
				Metadata: types.Metadata{
					Corpus: corpusID,
				},
			},
		})
	}
	return similarities, nil
}

func constructSimilaritiesFromResult(result *client.SearchResult) ([]*types.Similarity, error) {
	var idCol *entity.ColumnVarChar
	var documentIDCol *entity.ColumnVarChar
	var corpusIDCol *entity.ColumnVarChar
	//body, _ := json.Marshal(result)
	//fmt.Println(string(body))
	for _, field := range result.Fields {
		switch field.Name() {
		case FiledID:
			field.FieldData()
			if c, ok := field.(*entity.ColumnVarChar); ok {
				idCol = c

			}
		case FiledDocumentID:
			if c, ok := field.(*entity.ColumnVarChar); ok {
				documentIDCol = c
			}
		case FiledCorpus:
			if c, ok := field.(*entity.ColumnVarChar); ok {
				corpusIDCol = c
			}
		}
	}

	var similarities []*types.Similarity
	for i := 0; i < result.ResultCount; i++ {
		id, err := idCol.GetAsString(i)
		if err != nil {
			return nil, err
		}
		documentID, err := documentIDCol.GetAsString(i)
		if err != nil {
			return nil, err
		}

		corpusID, err := corpusIDCol.GetAsString(i)
		if err != nil {
			return nil, err
		}
		similarities = append(similarities, &types.Similarity{
			Chunk: &types.Chunk{
				ID:         id,
				DocumentID: documentID,
				Metadata: types.Metadata{
					Corpus: corpusID,
				},
			},
			Score: float64(result.Scores[i]),
		})
	}

	return similarities, nil
}
