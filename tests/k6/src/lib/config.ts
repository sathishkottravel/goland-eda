// Base URL of the api service. Override with -e BASE_URL=... or the K6_BASE_URL compose var.
export const BASE_URL = __ENV.BASE_URL || "http://api:8080";

// Default pass/fail thresholds shared by all tests.
export const defaultThresholds = {
  http_req_failed: ["rate<0.01"],
  http_req_duration: ["p(95)<500"],
};
