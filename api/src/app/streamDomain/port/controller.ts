import { Writable } from "stream"
import { ISocket } from "../../infraestructure/router"
import { IStreamPluginCommand } from "../useCase/pluginCommand"



export interface IStreamController {
  plugInHub(driveId: string, socket: ISocket, channel: string): { onClose: () => void }

}

export class StreamController implements IStreamController {

  private streamPluginCommand: IStreamPluginCommand
  constructor({ streamPluginCommand }: { streamPluginCommand: IStreamPluginCommand }) {
    this.streamPluginCommand = streamPluginCommand
  }

  public plugInHub(driveId: string, socket: ISocket, channel: string) {
    const { clientStream } = this.streamPluginCommand.plugInHub(driveId, socket.id)


    const writeStream = this.writeStreamOnObj((chunk) => {
      socket.emit(channel, chunk)
    })

    clientStream.pipe(writeStream)

    return {
      onClose: () => {
        console.info(`closing connection of ${socket.id}`)
      }
    }
  }

  private writeStreamOnObj(func: (chunk: any) => void) {
    return new Writable({
      objectMode: true,
      write(chunk, encoding, callback) {
        func(chunk)
        callback()
      }
    })
  }


}