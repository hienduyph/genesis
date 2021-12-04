import rlp from 'readline';

export function ask(text: string): Promise<string> {
  const rl = rlp.createInterface({
    input: process.stdin,
    output: process.stdout
  });
  return new Promise((resolve) => {
    rl.question(`${text}: `, (input) => {
      resolve(input);
      rl.close();
    });
  });
}