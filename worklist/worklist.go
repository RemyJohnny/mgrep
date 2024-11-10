package worklist

type Entry struct{
	Path string
}

type WorkList struct{
	jobs chan Entry
}

// Adds work to the channel
func (w *WorkList) Add(work Entry){
	w.jobs <- work
}

// Gets work out of the jobs channel
func (w *WorkList) Next() Entry{
	job := <- w.jobs
	return job
}

// creates a new workList with a buffered jobs channel
func New(bufSize int) WorkList{
	return WorkList{
		jobs: make(chan Entry,bufSize),
	}
}

// Creates a new job i.e Entry
func NewJob( path string) Entry{
	return Entry{path}
}

// spams all the workers with a blank job (Entry). Indication for them to terminate
func (w *WorkList) Finalize(workersNum int ){
	for i := 0; i < workersNum; i++{
		w.Add(Entry{""})
	}
}