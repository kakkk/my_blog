const path = require('path');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const ArcoWebpackPlugin = require('@arco-plugins/webpack-react');
module.exports = {
    entry: './src/index.tsx',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'bundle.js',
    },
    resolve: {
        extensions: ['.ts', '.tsx', '.js'],
    },
    module: {
        rules: [
            {
                test: /\.(ts|tsx)$/,
                exclude: /node_modules/,
                use: ['ts-loader'],
            },
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader']
            },
            {
                test: /\.less$/,
                use: [
                    'style-loader',
                    'css-loader',
                    'less-loader'
                ]
            }
        ],
    },
    devServer: {
        static: {
            directory: path.resolve(__dirname, 'dist'),
        },
        compress: true,
        port: 3000,
    },
    plugins: [
        new CopyWebpackPlugin({
            patterns: [
                { from: 'public/index.html', to: 'index.html' },
                { from: 'public/assets', to: 'assets' },
            ],
        }),
        new ArcoWebpackPlugin({}),
    ],
};
