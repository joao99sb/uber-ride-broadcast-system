import { Readable, Writable } from "stream"
import { QUEUE_NAME } from "../../../consts"
import { IQueue } from "../../services/queue/port"


export interface IHub {
  plugInHub(driveId: string, socketId: string, stream: Writable): void
}

type HubProps = {
  queue: IQueue
}

type QueueResponse = {
  order: string
  lat: string
  lng: string

}
export class Hub implements IHub {

  private queue: IQueue
  private queueStream: Readable
  private clients: Map<string, { stream: Writable, driveId: string }> = new Map
  constructor({ queue }: HubProps) {
    this.queue = queue
    this.queueStream = new Readable({
      objectMode: true,
      read() { }
    });
    this.clients = new Map()
    this.configQueue()
  }

  private configQueue() {


    this.queue.on('message', (data: QueueResponse) => {
      const { order, lat, lng } = data
      this.queueStream.push({ order, lat, lng })
    })

    this.queue.on("initiated", () => {
      this.queue.consumeQueue(QUEUE_NAME).catch(console.error)
    })
  }
  public startHub() {
    this.queueStream.pipe(this.broadcast())
  }
  public plugInHub(driveId: string, socketId: string, stream: Writable) {
    this.clients.set(socketId, { stream, driveId })
  }
  private broadcast() {

    return new Writable({
      objectMode: true,
      write: (chunk: QueueResponse, enc, cb) => {
        for (const [id, client] of this.clients) {
          const { stream, driveId } = client
          const { order } = chunk
          if (stream.writableEnded) {
            this.clients.delete(id)
            continue;
          }
          if (driveId === order) {
            stream.write(JSON.stringify(chunk))
          }
        }

        cb()
      }
    })
  }

}