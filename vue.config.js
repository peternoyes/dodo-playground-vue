module.exports = {
    configureWebpack: {
        module: {
            rules: [
                {
                test: /\.md$/,
                use: [
                    {
                        loader: "html-loader"
                    },
                    {
                        loader: "markdown-loader",
                        options: {
                            /* your options here */
                        }
                    }
                ]
                },
                {
                    test: /\.c$/,
                    use: "raw-loader"   
                },
                {
                    test: /\.wasm$/,
                    type: 'javascript/auto',
                    loaders: ['arraybuffer-loader'],
                }                        
            ],
        }
    }
  };