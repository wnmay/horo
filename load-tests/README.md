# How to run the test

cd load-tests

# Test Order service
.\run-test.ps1 order-service-test.js

# Test Chat service
.\run-test.ps1 chat-rabbit-latency-test.mjs
.\run-test.ps1 chat-service-test.js
