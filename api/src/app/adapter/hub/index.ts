import { Hub } from "../../hubDomain/port/controller";
import { queue } from "../../services/adapter";

export const hubCommand = new Hub({
  queue: queue
})