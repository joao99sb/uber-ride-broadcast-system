import { PassThrough, Writable } from "stream"
import { IHub } from "../../hubDomain/port/controller"


export interface IStreamPluginCommand {
  plugInHub(driveId: string, socketId: string): { clientStream: PassThrough }
}

export class StreamPluginCommand implements IStreamPluginCommand {

  private hubCommand: IHub
  constructor({ hubCommand }: { hubCommand: IHub }) {
    this.hubCommand = hubCommand
  }

  public plugInHub(driveId: string, socketId: any) {
    const { clientStream } = this.createClientStream()
    this.hubCommand.plugInHub(driveId, socketId, clientStream)

    return { clientStream }
  }

  private createClientStream() {
    const clientStream = new PassThrough(
      { objectMode: true }
    )

    return {
      clientStream
    }
  }

}