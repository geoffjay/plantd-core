package mdp

func popWorker(workers []*brokerWorker) (worker *brokerWorker, workers2 []*brokerWorker) {
	worker = workers[0]
	workers2 = workers[1:]
	return
}

func delWorker(workers []*brokerWorker, worker *brokerWorker) []*brokerWorker {
	for i := 0; i < len(workers); i++ {
		if workers[i] == worker {
			workers = append(workers[:i], workers[i+1:]...)
			i--
		}
	}
	return workers
}

func stringArrayToByte2D(in []string) (out [][]byte) {
	for _, str := range in {
		out = append(out, []byte(str))
	}
	return
}

func byte2DToStringArray(in [][]byte) (out []string) {
	for _, bytes := range in {
		out = append(out, string(bytes))
	}
	return
}
