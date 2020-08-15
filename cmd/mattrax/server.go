package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// Serve uses the arguments to create a HTTPS server that uses secure defaults and has gracefully shutdown support
func serve(addr string, domain string, httpsCertPath string, httpsKeyPath string, caCertPool *x509.CertPool, r http.Handler) {
	var srv = &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			NextProtos:               []string{"h2", "http/1.1"},
			// Mutual TLS
			ClientCAs: caCertPool,
			// ClientAuth:            tls.VerifyClientCertIfGiven, // tls.NoClientCert,
			// Standards from https://wiki.mozilla.org/Security/Server_Side_TLS
			MinVersion: tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.CurveP384,
				tls.CurveP521,
				tls.X25519,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServeTLS(httpsCertPath, httpsKeyPath); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server encountered an error")
		}
	}()
	log.Info().Str("addr", addr).Str("host", domain).Msg("Listening...")

	<-done
	log.Info().Msg("Finishing active connections. Please wait...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Failed to shutdown server")
	}
}
