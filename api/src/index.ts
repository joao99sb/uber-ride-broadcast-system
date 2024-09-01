import dotenv from 'dotenv'
dotenv.config()
import express from 'express'
import { createServer } from 'http'

import { MountRouterHandler } from './infra/router/express';

import { MountSocketHandler } from './infra/router/socketio';
import { streamController } from './app/adapter/stream/indext';
import { hubCommand } from './app/adapter/hub';


const PORT = process.env.PORT || 3001

async function main() {
  const app = express()
  const routerMounter = new MountRouterHandler({ app })
  routerMounter.config()

  const server = createServer(app)

  const socketHandler = new MountSocketHandler({ server, streamController })
  socketHandler.config()
  hubCommand.startHub()


  server.listen(PORT, () => {
    console.log(`Servidor rodando na porta ${PORT}`);
  })

}
main()




