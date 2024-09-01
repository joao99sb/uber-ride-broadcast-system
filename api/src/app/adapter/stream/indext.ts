import { StreamController } from "../../streamDomain/port/controller";
import { StreamPluginCommand } from "../../streamDomain/useCase/pluginCommand";
import { hubCommand } from "../hub";


const streamPluginCommand = new StreamPluginCommand({ hubCommand })
export const streamController = new StreamController({ streamPluginCommand }) 