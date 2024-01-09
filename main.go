package main

func main() {
	dataStream, done := GetDataStream()

	pipeline := NewPipeline(
		done,
		filterNegativeValues,
		filterSpecificValues,
		bufferValues,
	)

	consumer(done, pipeline.Run(dataStream))
}
