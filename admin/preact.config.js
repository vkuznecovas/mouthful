import asyncPlugin from 'preact-cli-plugin-async';

export default function (config, env, helpers) {
  if (process.env.HOMEPAGE) {
    config.output.publicPath = process.env.HOMEPAGE
  }
  asyncPlugin(config);
}