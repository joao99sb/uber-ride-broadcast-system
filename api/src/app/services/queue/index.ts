import amqp from 'amqplib'
import { EventEmitter } from 'stream';
import { QUEUE_HOST, QUEUE_NAME, QUEUE_PASS, QUEUE_PORT, QUEUE_USER, QUEUE_VHOST } from '../../../consts';
import { IQueue } from './port';




export class Queue extends EventEmitter implements IQueue {
  private url: string;
  private channel: amqp.Channel | null;
  private static instance: Queue
  constructor() {
    super();

    this.url = `amqp://${QUEUE_USER}:${QUEUE_PASS}@${QUEUE_HOST}:${QUEUE_PORT}${QUEUE_VHOST}`

    this.channel = null
    this.connect()
  }


  static getInstance() {
    if (!this.instance) {
      this.instance = new Queue()
    }
    return this.instance
  }

  addEvent(event: string, cb: (...args: any[]) => void): void {
    this.on(event, cb)
  }



  async consumeQueue(queue: string) {
    if (!this.channel) {
      throw new Error('Channel not initialized')
    }

    this.channel.consume(queue, (msg) => {
      if (msg !== null) {
        console.log(msg.content.toString())
        const content = JSON.parse(msg.content.toString())
        this.emit('message', content)
        this.channel!.ack(msg)
      }
    })
  }


  async connect() {
    try {
      const connection = await amqp.connect(this.url);
      const channel = await connection.createChannel()
      console.log('Connected to queue')
      this.channel = channel
      await this.assertQueue(channel, QUEUE_NAME)
      this.emit('initiated')
    } catch (err) {
      console.error(err);
    }
  }

  private async assertQueue(channel: amqp.Channel, queue: string, option?: amqp.Options.AssertQueue) {
    const defaultOption: amqp.Options.AssertQueue = {
      exclusive: false,
      durable: true,
      autoDelete: false,
      arguments: null
    }

    await channel.assertQueue(queue, { ...defaultOption, ...option })

  }

}