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
                    test: /\.wasm|bin$/,
                    type: 'javascript/auto',
                    loaders: ['arraybuffer-loader'],
                }                        
            ],
        }
    },
    chainWebpack: config => {
        config
            .plugin('html')
            .tap(args => {
                args[0].title = "Dodo Playground (BETA)";
                return args;
            })
    }    
  };