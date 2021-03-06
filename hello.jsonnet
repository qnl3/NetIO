{
    type: 'NetIOConfiguration',
    name: 'hello',
    description: 'hello world example configuration',
    version: 2,
    theme: 'dark',
    connections: [
        {
            name: 'helloworld',
            host: std.extVar('targetHost'),
            port: std.extVar('portNumber'),
        },
    ],
    pages: [
        {
            name: 'first',
            label: 'first',
            connection: 'helloworld',
            fitToScreen: true,
            items: [
                {
                    type: 'button',
                    width: 300,
                    height: 35,
                    left: 10,
                    top: 80,
                    shape: 'rounded',
                    icon: '',
                    label: 'Hello World!',
                    sends: [
                        'hello',
                    ],
                    background: '45,45,45',
                    border: '70,70,70',
                    textcolor: '230,230,230',
                },
            ],
            width: 320,
            height: 480,
            textcolor: '230,230,230',
            background: '40,40,40',
        },
        {
            name: 'second',
            label: 'sec',
            connection: 'helloworld',
            fitToScreen: true,
            items: [],
            width: 320,
            height: 480,
            textcolor: '230,230,230',
            background: '40,40,40',
        },
    ],
    device: 'iPhone',
    navigation: 'fix',
    pagebuttonwidth: 'static',
    style: 'flat',
    orientation: 'portrait',
    navBackground: '30,30,30',
    preventSleep: false,
    switchOnSwipe: true,
}
