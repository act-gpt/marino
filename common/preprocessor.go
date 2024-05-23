package common

import (
	"strings"
	"unicode/utf8"

	"github.com/act-gpt/marino/config/system"
	"github.com/act-gpt/marino/types"

	"github.com/go-aie/xslices"
)

type PreprocessorConfig struct {
	// ChunkTokenNum is the number of tokens for each text chunk.
	MaxTokens int

	// MinChunkCharNum is the minimum number of characters for each text chunk.
	MinTokens int

	Overlap int

	// MinChunkLenToEmbed is the minimum length in characters.
	// Chunks with shorter length will be discarded.
	//MinChunkLenToEmbed int

	// MaxChunkNum is the maximum number of chunks to generate from a text.
	//MaxChunkNum int

	// MaxChunkNum is the maximum number of chunks to generate from a text.

}

func (cfg *PreprocessorConfig) Init() *PreprocessorConfig {
	config := system.Config.Parser
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = config.MaxTokens
	}
	if cfg.MinTokens == 0 {
		cfg.MinTokens = config.MinTokens
	}
	if cfg.Overlap == 0 {
		cfg.Overlap = config.Overlap
	}
	/*
		if cfg.MinChunkLenToEmbed == 0 {
			cfg.MinChunkLenToEmbed = config.MinChunkLenToEmbed
		}
		if cfg.MaxChunkNum == 0 {
			cfg.MaxChunkNum = config.MaxChunkNum
		}
	*/
	return cfg
}

// Preprocessor splits a list of documents into chunks.
type Preprocessor struct {
	encoder *dummyTokenizer
	cfg     *PreprocessorConfig
}

func NewPreprocessor(cfg *PreprocessorConfig) *Preprocessor {
	return &Preprocessor{
		encoder: Encoder,
		cfg:     cfg.Init(),
	}
}

func (p *Preprocessor) Preprocess(docs ...*types.Document) (map[string][]*types.Chunk, error) {
	chunkMap := make(map[string][]*types.Chunk)
	for _, doc := range docs {
		docID := doc.ID
		meta := doc.Metadata
		if docID == "" {
			docID = GetUUID()
		}
		textChunks, err := p.split(doc.Text)
		if err != nil {
			return nil, err
		}
		for _, text := range textChunks {
			chunkMap[docID] = append(chunkMap[docID], &types.Chunk{
				ID:         GetUUID(),
				DocumentID: docID,
				Metadata:   meta,
				Text:       text,
			})
		}
	}

	return chunkMap, nil
}

/*
# Constants
CHUNK_SIZE = 200  # The target size of each text chunk in tokens
MIN_CHUNK_SIZE_CHARS = 350  # The minimum size of each text chunk in characters
MIN_CHUNK_LENGTH_TO_EMBED = 5  # Discard chunks shorter than this
MAX_NUM_CHUNKS = 10000  # The maximum number of chunks to generate from a text
*/

// split converts the text into chunks.
//
// The splitting algorithm is borrowed from https://github.com/openai/chatgpt-retrieval-plugin/blob/88d983585816b7f298edb0cabf7502c5ccff370d/services/chunks.py#L22-L96.
func (p *Preprocessor) split(text string) ([]string, error) {

	if text == "" || strings.TrimSpace(text) == "" {
		return nil, nil
	}

	// Convert the document text into runes.
	runes := []rune(text)

	var chunks []string

	var i int

	for i < len(runes)-1 {
		// Take the first ChunkTokenNum tokens as a chunk.
		chunkRunes, err := p.encoder.Encode(runes[i:], p.cfg.MaxTokens)
		if err != nil {
			return nil, nil
		}

		// Skip the chunk if it is empty or whitespace.
		chunkText := string(chunkRunes)

		if strings.TrimSpace(chunkText) == "" {
			i += len(chunkRunes)
			continue
		}

		// Find the last period or punctuation mark in the chunk.)
		// Note that here we count the index in runes.
		var lastPuncIdx = -1
		for _, punc := range []rune{
			'.', '?', '!',
			'。', '？', '！', '；',
			'\n',
		} {
			lastPuncIdx = xslices.Max(lastPuncIdx, lastRuneIndex(chunkText, punc))
		}
		if lastPuncIdx != -1 && lastPuncIdx > p.cfg.MinTokens {
			// Truncate the chunk text at the punctuation mark.
			chunkText = string([]rune(chunkText)[:lastPuncIdx+1])
		}
		/*
			if len(runes)-i < p.cfg.ChunkOverlap {
				val := chunks[len(chunks)-1]
				fmt.Println(val + string(chunkRunes))
			}
		*/
		// 把换行符都去掉？
		trimmedChunkText := strings.TrimSpace(strings.ReplaceAll(chunkText, "\n", " "))
		//trimmedChunkText := strings.TrimSpace(chunkText)
		if utf8.RuneCountInString(trimmedChunkText) > p.cfg.MinTokens {
			chunks = append(chunks, trimmedChunkText)
		}
		if i > len(runes)-p.cfg.MinTokens {
			i += utf8.RuneCountInString(chunkText)
		} else {
			i += utf8.RuneCountInString(chunkText) - p.cfg.Overlap
		}
		//chunkNum += 1
	}

	// Handle the remaining runes.
	if i < len(runes) {
		remainingText := string(runes[i:])

		//trimmedRemainingText := strings.TrimSpace(strings.ReplaceAll(remainingText, "\n", " "))
		trimmedRemainingText := strings.TrimSpace(remainingText)
		if utf8.RuneCountInString(trimmedRemainingText) > p.cfg.MinTokens {
			chunks = append(chunks, trimmedRemainingText)
		}
	}

	return chunks, nil
}

func lastRuneIndex(s string, r rune) int {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == r {
			return i
		}
	}
	return -1
}
