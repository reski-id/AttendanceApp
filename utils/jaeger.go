package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func JaegerTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize the tracer
		tracer := opentracing.GlobalTracer()
		spanCtx, err := tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Printf("Failed to extract span context: %v", err)
		}

		span := tracer.StartSpan(
			c.Request.URL.Path,
			opentracing.ChildOf(spanCtx),
			opentracing.Tag{Key: "http.method", Value: c.Request.Method},
			opentracing.Tag{Key: "http.url", Value: c.Request.URL.Path},
		)
		defer span.Finish()

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
		c.Next()
	}
}
