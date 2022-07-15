### Q. Out of above two backends, performing the same operations but using different backend language, which one is more efficient, performant, scalable for future? Also provide evidences to prove your point.

- Golang API is more efficient, performant and scalable for the future. Go is a statically-typed language. This helps in reducing a lot of production bugs and issues, whereas similar functionality in NodeJS (dynamic typing) requires extra setup and tooling for NodeJS API. Many bugs can be caught on compile-time itself, and facilitates the validation of data types of our models.

- Go is compiled and the executable is then deployed on the development/production servers. This is always preferable to interpreted languages like Javascript. The V8 engine that runs NodeJS application comes very close, as it does Just-in-Time compilation - which tries to prevent long compilation times.

- Go is good for CPU intensive tasks, but NodeJS is preferred for IO intensive tasks.

- For the NodeJS case, as requests are processed and events are triggered, messages are queued along with their callback functions. The queue is polled for the next message, and when a message is seen, callback for that message is executed. There are issues when there is high traffic, and the API is contacting other 3rd-party APIs, there might be timeouts, and events/messages/callbacks are placed on overloaded queue, and the callbacks are processed one by one. The single-threaded model has its disadvantages, which is why NodeJS has introduced worker threads recently to tackle some of these issues.

- Go required explicit error checking which can lead to carefully-designed, flawless apps, but NodeJS uses the more traditional throw-and-catch model for error handling.

- Go's strength lies in in-built & streamlined concurrency, also the optimized garbage collection and memory management. NodeJS has a non-blocking asynchronous system which makes it seem like multi-threaded but it is not actually that.

- Golang comes with goroutines which are like lightweight threads, and channels can be used for inter-goroutine communication. These concurrency primitives are built into the language, which help to develop for use-cases that have heavy processing going on in the backend.

- Both have their pros and cons, and to know which is better for the company's requirements, we can do load testing and stress testing. Grafana's K6 is a modern tool to do such kind of testing and also commit these tests in the repository to make sure they run in a CI/CD setup post-development. There are other benchmarking tools, but comparisons need to be made carefully, since there might be bottlenecks in the type of processing going on for each API route.

- To summarize, Golang is scalable due to its concurrency, efficient due to its static typing and compile-time error catching, & performant since it is really close to machine-level code like C++ or Java. But it depends on the scale and the requirements, where NodeJS tooling and open-source packages can help in quickly setting up prototypes and proofs-of-concept.