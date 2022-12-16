package chunks

import (
	"gorm.io/gorm"
	"server/conf"
	"server/db/models/blocks"
)

type Chunk struct {
	gorm.Model

	// refers to the top-left coord of a chunk
	X int `json:"x"`
	Y int `json:"y"`

	Blocks []blocks.Block `json:"blocks"`
}

func ToChunkCoords(pos []int) []int { // converts block coordinates to chunk coordinates
	chunkCoords := make([]int, 2)
	for i := range pos {
		chunkCoords[i] = int(pos[i] / conf.CHUNK_SIZE)
		if chunkCoords[i] < 0 {
			chunkCoords[i] -= 1
		}
	}

	return chunkCoords
}

func ChunksInRenderDist(curCunk []int) [][]int { // returns all chunks within render distance
	chunks := [][]int{}
	i := 0
	for x := curCunk[0] - conf.RENDER_DISTANCE; x < curCunk[0]+conf.RENDER_DISTANCE; x++ {
		for y := curCunk[1] - conf.RENDER_DISTANCE; y < curCunk[1]+conf.RENDER_DISTANCE; y++ {
			chunks = append(chunks, []int{x, y})
			i++
		}
	}
	return chunks
}

func ChunkSpanCoords(curChunk []int) (TL [2]int, BR [2]int) { // returns the top-left and bottom-right chunks in a span
	/*
	 this is done by sending the top left(TL) and bottom right(BR) chunk coordinates
	 all other chunks can be found within
	*/
	TL = [2]int{curChunk[0] - conf.RENDER_DISTANCE, curChunk[1] - conf.RENDER_DISTANCE}
	BR = [2]int{curChunk[0] + conf.RENDER_DISTANCE, curChunk[1] + conf.RENDER_DISTANCE}
	return TL, BR
}

func ChunkXSpan(L int, R int) []int {
	/*
		Returns all X coordinates in a given chunk span
		takes in the Left and Right most coordinates and returns the range between them
		used for our SQL queries
	*/
	Xspan := []int{}
	for i := L; i < R; i++ {
		Xspan = append(Xspan, i)
	}
	return Xspan
}

func ChunkYSpan(T int, B int) []int {
	/*
		Returns all Y coords in a given chunks span using the top and bottom most coords and returns a range between em
		used for SQL queries
	*/
	Yspan := []int{}
	for i := T; i < B; i++ {
		Yspan = append(Yspan, i)
	}
	return Yspan
}

func ChunkSpan(curChunk []int) (TL [2]int, BR [2]int, xspan []int, yspan []int) {
	/*
		Wrapper for ChunkSpanCoords, ChunkXSpan, ChunkYSpan
		given a set of chunk coordinates, this determines the span of chunks to be rendered.
		Returns the top-left and bottom-right coordinates(TL, BR) and the span of X and Y values within that chunk span.

		SQL queries will require the X and Y spans and the client can use the TL and BR coordinates to generate all chunks.
		This(2C ~= Xspan + Yspan + TL+ BR) is much more efficient than sending every possible chunk combination(C^2).
	*/

	TL, BR = ChunkSpanCoords(curChunk)
	xspan = ChunkXSpan(TL[0], BR[0])
	yspan = ChunkYSpan(TL[1], BR[1])

	return TL, BR, xspan, yspan
}