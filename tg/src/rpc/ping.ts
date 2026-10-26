import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { PingService } from "../proto/tg_ping_pb";

console.log("BASE", process.env.SERVER_BASE_URL)

const transport = createConnectTransport({
  baseUrl: process.env.SERVER_BASE_URL || "http://localhost:8080",
	httpVersion: "1.1",
});

export const pingRPC = createClient(PingService, transport);
