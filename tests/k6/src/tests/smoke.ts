import { check, sleep } from "k6";
import type { Options } from "k6/options";
import { defaultThresholds } from "../lib/config";
import { get, graphql } from "../lib/http";

export const options: Options = {
  vus: 1,
  duration: "10s",
  thresholds: { ...defaultThresholds, checks: ["rate==1"] },
};

export default function () {
  const health = get("/health");
  check(health, {
    "health 200": (r) => r.status === 200,
    "health ok": (r) => r.json("status") === "ok",
  });

  const hello = get("/api/v1/hello?name=k6");
  check(hello, {
    "hello 200": (r) => r.status === 200,
    "hello message": (r) => r.json("message") === "Hello, k6!",
  });

  const gql = graphql(`query ($name: String) { hello(name: $name) }`, { name: "k6" });
  check(gql, {
    "graphql 200": (r) => r.status === 200,
    "graphql hello": (r) => r.json("data.hello") === "Hello, k6!",
  });

  const playground = get("/playground");
  check(playground, { "playground 200": (r) => r.status === 200 });

  sleep(1);
}
