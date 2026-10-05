import http from "k6/http";
import { BASE_URL } from "./config";

export function get(path: string) {
  return http.get(`${BASE_URL}${path}`, { tags: { name: path } });
}

export function graphql(query: string, variables: Record<string, unknown> = {}) {
  return http.post(`${BASE_URL}/graphql`, JSON.stringify({ query, variables }), {
    headers: { "Content-Type": "application/json" },
    tags: { name: "/graphql" },
  });
}
