import { Server, Socket } from "socket.io"
import { DefaultEventsMap } from "socket.io/dist/typed-events"
import { IStreamController } from "../../../app/streamDomain/port/controller"


type ConfigParams = {
  server: any
  streamController: IStreamController
}

export type ISocketIO = Socket<DefaultEventsMap, DefaultEventsMap, DefaultEventsMap, any>

export class MountSocketHandler {


  private io: Server
  private server: any
  private streamController: IStreamController
  constructor({ server, streamController }: ConfigParams) {
    this.server = server
    this.io = new Server(this.server, {
      cors: {
        origin: "*",
      }
    })
    this.streamController = streamController
  }

  public config() {


    this.io.on('connection', (socket) => {
      console.log('a user connected');
      this.mountSocket(socket)

    });


  }


  public mountSocket(socket: Socket<DefaultEventsMap, DefaultEventsMap, DefaultEventsMap, any>) {
    socket.on('travel_id', async (driveId: any) => {
      const { onClose } = this.streamController.plugInHub(driveId, socket, 'travel_id');
      socket.on('disconnect', onClose)

    })


    socket.on('disconnect', () => {
      console.log('user disconnected');
    });
  }



}



