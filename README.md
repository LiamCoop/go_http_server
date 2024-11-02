Core features
- [ ] Keep Alive
- [ ] Timeouts
- [ ] Max Requests
- [ ] Routing: Create a mechanism to route requests to specific handlers based on the path and method. For example, a simple map could direct GET /api/v1/resource to a specific handler function. 
- [ ] Middleware Support (optional but helpful): Adding middleware can let you handle things like logging, error handling, and authentication before the request reaches your main handler.- 
- [ ] Error Handling: Implement proper error responses, such as 404 Not Found or 500 Internal Server Error, when requests cannot be fulfilled.
- [ ] Timeouts and Connection Closing: Now that you’re looking at keep-alive functionality, add logic to close idle connections after a set period. Use context timeouts to handle long-running requests gracefully.
- [ ] Chunked Transfer Encoding (for larger responses): This is useful if you want to handle large or streaming data responses in a way that doesn’t require knowing the full content length upfront.

More advanced features

- [ ] TLS/HTTPS Support: Securing the connection is vital for real-world applications. Setting up HTTPS with Go’s crypto/tls package would add robustness to your server.
- [ ] Compression Support: Implement gzip/deflate support to reduce data transfer size, especially for large payloads.
- [ ] Cookie Handling: Implementing support for Set-Cookie and parsing Cookie headers can help prepare your server for session-based or authenticated interactions.

