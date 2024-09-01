import { Express } from 'express'
import cors from 'cors'
import express from 'express'


export interface IMountRouterHandler {
  config({ app }: { app: Express }): void
}

export class MountRouterHandler {

  private app: Express

  constructor({ app }: { app: Express }) {
    this.app = app
  }



  public config() {
    const app = this.app
    app.use(cors({
      origin: "*",
      methods: ["GET", "POST"]
    }))

    app.use(express.json())

    this.mountRoutes()
  }

  public mountRoutes() {
    const app = this.app

    app.route('/health-check').get(async (req, res) => {
      res.status(200).send({
        status: 'UP',
      });
    });
  }
}