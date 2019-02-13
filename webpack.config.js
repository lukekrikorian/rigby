const path = require("path");

module.exports = {
    mode: "production",
    entry: ["./client/index.jsx"],
    output: {
        filename: "bundle.js",
        path: path.resolve(__dirname + "/static")
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: ['cache-loader',
                {
                    loader: 'babel-loader',
                    options: {
                        babelrc: false,
                        presets: ["@babel/env", "@babel/react"]
                    }
                }]
            },
            {
                test: /\.css$/,
                use: ["style-loader", "css-loader"]
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