import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { PingService } from "../proto/api_ping_pb";

const transport = createConnectTransport({
	baseUrl: import.meta.env["VITE_SERVER_BASE_URL"] || "http://localhost:9000",
});

export const api = createClient(PingService, transport);
