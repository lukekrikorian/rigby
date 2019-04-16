const path = require("path");

module.exports = {
    entry: "./client/index.jsx",
    output: {
        filename: "[name].[contentHash].js",
        path: path.resolve(__dirname + "/static/dist")
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: ["cache-loader",
                {
                    loader: 'babel-loader',
                    options: {
                        babelrc: false,
                        presets: ["@babel/env", "@babel/react"]
                    }
                }]
            }
        ]
    },
    resolve: {
        alias: {
            "react": "preact-compat",
            "react-dom": "preact-compat"
        },
        extensions: [".js", ".jsx"]
    }
};