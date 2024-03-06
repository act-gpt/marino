package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/config/system"

	"github.com/pgvector/pgvector-go"
)

type Segment struct {
	Id          string    `json:"id"`
	KnowledgeId string    `json:"knowledge_id"`
	Corpus      string    `json:"corpus"`
	Index       int       `json:"index"`
	Text        string    `json:"text"`
	Source      string    `json:"source"`
	Url         string    `json:"url"`
	Sha         string    `json:"sha"`
	Embedding   []float32 `json:"embedding" gorm:"-:all"`
	Score       float64   `json:"score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (segment *Segment) Insert() error {
	var err error
	if segment.Id == "" {
		segment.Id = common.GetUUID()
	}
	sha := segment.Sha
	if sha == "" {
		sha = common.ContentSha(segment.Text)
	}
	now := time.Now()
	err = DB.Exec("INSERT INTO segments (id, knowledge_id, corpus, index, text, source, url, sha, embedding, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", segment.Id, segment.KnowledgeId, segment.Corpus, segment.Index, segment.Text, segment.Source, segment.Url, sha, pgvector.NewVector(segment.Embedding), now, now).Error
	return err
}

func (segment *Segment) Update() error {
	var err error
	sha := segment.Sha
	if sha == "" {
		sha = common.ContentSha(segment.Text)
	}
	now := time.Now()
	err = DB.Exec("UPDATE segments SET knowledge_id=$1, corpus=$2, index=$3, text=$4, source=$5, url=$6, sha=$7, embedding=$8, updated_at=$9 WHERE id = $10", segment.KnowledgeId, segment.Corpus, segment.Index, segment.Text, segment.Source, segment.Url, sha, pgvector.NewVector(segment.Embedding), now, segment.Id).Error
	return err
}

func (segment *Segment) Delete() error {
	if segment.Id == "" {
		return fmt.Errorf("id 为空！")
	}
	err := DB.Delete(segment).Error
	return err
}

func InitSegments(conf system.SystemConfig) {
	c := conf.Db
	dimension := c.Dimension
	if dimension == 0 {
		dimension = 768
	}
	// for vector embedding
	DB.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	sql := `CREATE TABLE IF NOT EXISTS segments (
			id VARCHAR(191) NOT NULL,
			knowledge_id VARCHAR(32) DEFAULT NULL,
			corpus VARCHAR,
			index BIGINT DEFAULT 0,
			text TEXT,
			source VARCHAR,
			url VARCHAR,
			sha VARCHAR,
			embedding vector(` + strconv.Itoa(dimension) + `),
			created_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
			updated_at TIMESTAMP(3) WITHOUT TIME ZONE DEFAULT NULL,
			PRIMARY KEY (id)
		  );`
	DB.Exec(sql)
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_segments_knowledge_id ON segments (knowledge_id);")
	/*
		if c.IndexType == "hnsw" {
			DB.Exec("CREATE INDEX IF NOT EXISTS idx_segments_embedding ON segments USING hnsw (embedding vector_ip_ops);")
		} else {
			DB.Exec("CREATE INDEX IF NOT EXISTS idx_segments_embedding ON segments USING ivfflat (embedding vector_ip_ops) WITH (lists = 100);")
		}
	*/
	if c.IndexType == "hnsw" {
		DB.Exec("CREATE INDEX IF NOT EXISTS idx_segments_embedding ON segments USING hnsw(embedding vector_cosine_ops) WITH (m = 24, ef_construction = 100)")
	} else {
		DB.Exec("CREATE INDEX IF NOT EXISTS idx_segments_embedding ON segments USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100)")
	}
}
func QueryEmbedding(embedding []float32, corpus string, limit int, score float64) ([]Segment, error) {
	ebd := pgvector.NewVector(embedding)
	var result []Segment
	//err := DB.Raw("SELECT id, knowledge_id, corpus, index, text, sha, (embedding <#> $1) * -1 AS score, created_at, updated_at FROM segments WHERE (embedding <#> $1) * -1 >= $3 AND corpus = $4 ORDER BY score DESC LIMIT $2", ebd, limit, score, corpus).Scan(&result).Error
	err := DB.Raw("SELECT id, knowledge_id, corpus, index, text, sha, 1 - (embedding <=> $1) AS score, created_at, updated_at FROM segments WHERE corpus = $3 ORDER BY score DESC LIMIT $2", ebd, limit, corpus).Scan(&result).Error
	fmt.Println(result)
	return result, err
}

func FindSegment(id string) (Segment, error) {
	var result Segment
	err := DB.Raw("SELECT id, knowledge_id, corpus, index, text, source, url, sha, created_at, updated_at FROM segments WHERE id = ? LIMIT 1", id).Scan(&result).Error
	return result, err
}

func FindSegments(ids []string) ([]Segment, error) {
	var result []Segment
	err := DB.Raw("SELECT id, knowledge_id, corpus, index, text, source, url, sha, created_at, updated_at FROM segments WHERE id IN ?", ids).Scan(&result).Error
	return result, err
}

func DeleteSegments(ids []string) error {
	err := DB.Exec("DELETE FROM segments WHERE knowledge_id IN ?", ids).Error
	return err
}
