import { EventEmitter } from "stream"

export interface IQueue extends EventEmitter {
  connect(): Promise<void>
  consumeQueue(queue: string): Promise<void>
  addEvent(event: string, cb: (...args: any[]) => void): void
}