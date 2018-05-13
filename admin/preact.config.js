export default function (config, env, helpers) {
    if (process.env.HOMEPAGE) {
        config.output.publicPath = process.env.HOMEPAGE        
    }
}