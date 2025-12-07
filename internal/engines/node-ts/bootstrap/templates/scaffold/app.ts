import { start } from "./vendor/server";

async function main() {
  await start();
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
