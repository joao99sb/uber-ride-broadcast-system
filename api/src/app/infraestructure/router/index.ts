
export interface ISocket {
  id: string
  emit(channel: string, message: any): void

} 