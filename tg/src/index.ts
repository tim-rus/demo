import { errAsync } from "./lib/errors"
import { pingRPC } from "./rpc/ping"

console.log("hi from tg")

async function main() {
  const [res, err] = await errAsync(pingRPC.ping({
    msg: "hi from tg",
  }))

  if (err) {
    console.error("failed to ping", err)
  }

  console.log("ping", res?.response, res?.timestamp?.seconds)
}

setTimeout(()=>main(), 5000)
