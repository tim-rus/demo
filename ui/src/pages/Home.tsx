import { Button } from "@/components/ui/button"
import { errAsync } from "@/lib/errors"
import { api } from "@/rpc/ping"
import { useState } from "react"

export default function Home() {
  const [state, setState] = useState("waiting...")

  const click = async () => {
    const [res, err] = await errAsync(api.ping({msg:"hi from ui"}))
    if (err) {
      console.error("failed to up", err)
      setState("error")
      return
    }
    setState(`${res.response}: ${res.timestamp.seconds}`)
  }
  return (
    <div>
      home: {import.meta.env.VITE_SERVER_BASE_URL}
      <div>
        State: {state}
        <Button onClick={click}>Go!</Button>
      </div>
    </div>
  );
}
