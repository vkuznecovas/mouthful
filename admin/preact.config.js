import asyncPlugin from 'preact-cli-plugin-async';
import envVars from 'preact-cli-plugin-env-vars';
import * as configFile from '../config.json';

export default (config, env, helpers) => {
  if (process.env.HOMEPAGE) {
    config.output.publicPath = process.env.HOMEPAGE;
  }

  process.env.PREACT_APP_URL = configFile.api.bindAddress;
  process.env.PREACT_APP_PORT = configFile.api.port;

  asyncPlugin(config);
  envVars(config, env, helpers);
}