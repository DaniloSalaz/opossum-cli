#!/usr/bin/env node

const yargs = require('yargs');
const axios = require('axios').default;
const BASE_URL = 'https://api.abuseipdb.com/api/v2/';

const argv = yargs
    .command('check', 'Retorna los detalles de una IP address', {
        ipAddress: {
            description: 'IP a verificar. (Required)',
            alias: 'ip',
            type: 'string',
        },
        maxDays: {
            description: 'Determina el tiempo (en días) máximo del reporte. Valor por defecto 30 días, mínimo (1) máximo (365)',
            alias: 'd',
            type: 'number',
        },
        verbose: {
            description: 'Verbose, retorna la toda la infomación.',
            alias: 'v',
            type: 'boolean',
        }
    }).command('blacklist', 'Retorna la lista negra de todas las IPs reportadas', {
        min: {
            description: 'Confidence minimum. Valor por defecto 100, mínumo(25) máximo (100)',
            alias: 'm',
            type: 'number',
        },
        limit: {
            description: 'Límite de IPs. Valor por defecto 10.000',
            alias: 'l',
            type: 'number',
        },
        plaintext: {
            description: 'Response en texto plano',
            alias: 'text',
            type: 'boolean'

        },
        desc: {
            description: 'Ordenar descendente',
            type: 'boolean'
        },
        last: {
            description: 'Ordenar descendente las IPs mas recientes',
            type: 'boolean'
        }
    })
    .option('key', {
        alias: 'k',
        description: 'ApiKey (Required)',
        type: 'string',
    })
    .help()
    .alias('help', 'h')
    .argv;


function run(command, apiKey, params) {
    axios.defaults.baseURL = BASE_URL;
    axios.defaults.headers['key'] = apiKey;

    let action = '';

    if (command === 'check') action = 'check';
    if (command === 'blacklist') action = 'blacklist';

    axios.get(action, params).then(res => {
        console.log(res.data);
    }).catch(error => {
        console.log(error.response.data);
    })

};
function main() {
    if (argv._.includes('check')) {
        if (argv.ipAddress || argv.key) {
            let options = {
                params: {
                    'ipAddress': argv.ipAddress,
                    'maxAgeInDays': argv.maxDays || 30,
                }
            }
            if (argv.verbose) options.params.verbose = true;
            const apiKey = argv.key;
            run('check', apiKey, options);

        } else {
            yargs.showHelp();
        }
    }else if (argv._.includes('blacklist')) {
        if (argv.key) {
            const options = {
                params: {
                    'confidenceMinimum': argv.min || 100,
                    'limit': argv.limit || 10000,
                }
            }

            if (argv.desc) options.params.abuseConfidenceScore = true;
            if (argv.last) options.params.lastReportedAt = true;
            if (argv.plaintext) options.params.plaintext = true;

            const apiKey = argv.key;
            run('blacklist', apiKey, options);

        } else {
            yargs.showHelp();
        }
    }else {
        yargs.showHelp();
    }

}

main();