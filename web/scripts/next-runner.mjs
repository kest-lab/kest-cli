import { spawn } from 'node:child_process';
import { createRequire } from 'node:module';

const [command, ...restArgs] = process.argv.slice(2);

if (!command) {
  process.exit(1);
}

const nextArgs = [command, ...restArgs];
const port = process.env.PORT?.trim();
const hostname = process.env.HOSTNAME?.trim();

if (port) {
  nextArgs.push('--port', port);
}

if (hostname) {
  nextArgs.push('--hostname', hostname);
}

const require = createRequire(import.meta.url);
const nextBin = require.resolve('next/dist/bin/next');

const child = spawn(process.execPath, [nextBin, ...nextArgs], {
  stdio: 'inherit',
  env: process.env,
});

child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }

  process.exit(code ?? 1);
});
