package splitters

import (
	"sort"
	"strings"

	"github.com/act-gpt/marino/common"
	"github.com/act-gpt/marino/model"
	"github.com/act-gpt/marino/types"
)

type Header struct {
	Level int
	Name  string
	Data  string
}

type Line struct {
	Metadata map[string]string
	Content  string
}

type Document struct {
	Metadata map[string]string
	Content  string
}

type MarkdownHeaderTextSplitter struct {
	HeadersToSplitOn []string
	ReturnEachLine   bool
	StripHeaders     bool
}

func NewMarkdownHeaderTextSplitter(headersToSplitOn []string, returnEachLine bool, stripHeaders bool) *MarkdownHeaderTextSplitter {
	sort.Slice(headersToSplitOn, func(i, j int) bool {
		return len(headersToSplitOn[i]) > len(headersToSplitOn[j])
	})
	return &MarkdownHeaderTextSplitter{HeadersToSplitOn: headersToSplitOn, ReturnEachLine: returnEachLine, StripHeaders: stripHeaders}
}

func (m *MarkdownHeaderTextSplitter) AggregateLinesToChunks(lines []Line) []Document {
	var aggregatedChunks []Line

	for _, line := range lines {
		if len(aggregatedChunks) > 0 {
			prev := aggregatedChunks[len(aggregatedChunks)-1]
			if prev.Metadata != nil {
				if prev.Metadata["Data"] == line.Metadata["Data"] {
					aggregatedChunks[len(aggregatedChunks)-1].Content += "  \n" + line.Content
				}
			}
		} else {
			aggregatedChunks = append(aggregatedChunks, line)
		}
	}

	var docs []Document
	for _, chunk := range aggregatedChunks {
		docs = append(docs, Document{Metadata: chunk.Metadata, Content: chunk.Content})
	}
	return docs
}

func (m *MarkdownHeaderTextSplitter) SplitText(text string) []Document {
	lines := strings.Split(text, "\n")
	var linesWithMetadata []Line
	var currentContent []string
	currentMetadata := make(map[string]string)
	var headerStack []Header

	for _, line := range lines {
		strippedLine := strings.TrimSpace(line)
		loop := false
		for _, sep := range m.HeadersToSplitOn {
			if strings.HasPrefix(strippedLine, sep) {
				currentHeaderLevel := strings.Count(sep, "#")
				for len(headerStack) > 0 && headerStack[len(headerStack)-1].Level >= currentHeaderLevel {
					headerStack = headerStack[:len(headerStack)-1]
				}

				header := Header{Level: currentHeaderLevel, Name: sep, Data: strings.TrimSpace(strippedLine[len(sep):])}
				headerStack = append(headerStack, header)
				currentMetadata[sep] = header.Data

				if len(currentContent) > 0 {
					linesWithMetadata = append(linesWithMetadata, Line{Metadata: currentMetadata, Content: strings.Join(currentContent, "\n")})
					currentContent = nil
				}

				if !m.StripHeaders {
					currentContent = append(currentContent, strippedLine)
				}
				loop = true
				break
			}
		}
		if !loop {
			if len(currentContent) > 0 && len(strippedLine) == 0 {
				linesWithMetadata = append(linesWithMetadata, Line{Metadata: currentMetadata, Content: strings.Join(currentContent, "\n")})
				currentContent = nil
			} else if len(strippedLine) > 0 {
				currentContent = append(currentContent, strippedLine)
			}
		}
		loop = false
	}

	if len(currentContent) > 0 {
		linesWithMetadata = append(linesWithMetadata, Line{Metadata: currentMetadata, Content: strings.Join(currentContent, "\n")})
	}

	return m.AggregateLinesToChunks(linesWithMetadata)
}

type MdPreprocess struct {
	cfg *PreprocessorConfig
}

func MdPreprocessor(cfg *PreprocessorConfig) *MdPreprocess {
	return &MdPreprocess{
		cfg: cfg.Init(),
	}
}

func (p *MdPreprocess) Preprocess(doc types.Document, bot model.BotSetting) (map[string][]*types.Chunk, []string, error) {
	chunkMap := make(map[string][]*types.Chunk)

	docID := doc.ID
	meta := doc.Metadata
	if docID == "" {
		docID = common.GetUUID()
	}
	textChunks, err := p.split(doc.Text)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Println("Chuncks size: ", len(textChunks))
	for _, textChunk := range textChunks {
		id := common.GetUUID()
		// return chunks
		chunkMap[docID] = append(chunkMap[docID], &types.Chunk{
			ID:         id,
			DocumentID: docID,
			Text:       textChunk,
			Metadata:   meta,
		})
	}

	return chunkMap, []string{}, nil
}

// spilt text to chunks
func (p *MdPreprocess) split(text string) ([]string, error) {
	headersToSplitOn := []string{"#", "##", "###"}
	markdownTextSplitter := MarkdownHeaderTextSplitter{
		HeadersToSplitOn: headersToSplitOn,
	}
	mdHeaderSplits := markdownTextSplitter.SplitText(text)
	lst := []string{}
	for _, doc := range mdHeaderSplits {
		lst = append(lst, doc.Content)
	}
	return lst, nil
}
