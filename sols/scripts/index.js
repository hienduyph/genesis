async function main() {
  const accounts = await ethers.provider.listAccounts();
  console.log(accounts);
}

main()
  .then(() => process.exit(0))
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
