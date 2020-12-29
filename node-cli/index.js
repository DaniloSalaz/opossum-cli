#!/usr/bin/env node

const yargs = require('yargs');
const axios = require('axios').default;
const BASE_URL = 'https://api.abuseipdb.com/api/v2/';

const argv = yargs
    .command('abuseip', "Requests to https://www.abuseipdb.com/ ", (yargs) => {
        yargs
        .usage('usage: $0 abuseip <item>')
        .command('check', 'IP Address report', {
            key: {
                alias: 'k',
                description: 'ApiKey (Required)',
                type: 'string',
            },
            ipAddress: {
                description: 'IP to verify. (Required)',
                alias: 'ip',
                type: 'string',
            },
            maxDays: {
                description: 'Determines how far back in time(days). The default is 30, min (1) max (365)',
                alias: 'd',
                type: 'number',
            },
            verbose: {
                description: 'Verbose',
                alias: 'v',
                type: 'boolean',
            }
        })
        .command('blacklist', 'Blacklist of reported IPs', {
            key: {
                alias: 'k',
                description: 'ApiKey (Required)',
                type: 'string',
            },
            min: {
                description: 'Confidence minimum. The default is 30, min (25) max (100)',
                alias: 'm',
                type: 'number',
            },
            limit: {
                description: 'Limit return IPs. The default 10.000',
                alias: 'l',
                type: 'number',
            },
            plaintext: {
                description: 'Response plain text',
                alias: 'text',
                type: 'boolean'
    
            },
            desc: {
                description: 'Descending order',
                type: 'boolean'
            },
            last: {
                description: 'Descending order of the last IPs',
                type: 'boolean'
            }
        })
        .updateStrings({
            'Commands:': 'item:'
        })
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