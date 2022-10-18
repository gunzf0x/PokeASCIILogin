# Installation
I provide 3 different ways to install this script:
#### I. If you use Zsh or Bash, I recommend to install it with a simple Bash script
1. Download `install.sh` script from the repo using `wget` or `curl` (we use `wget` in this example)
  
    ```
     wget https://raw.githubusercontent.com/p4nchit0z/PokeASCIILogin/main/install.sh -O install.sh
     ```

2. Execute the script
   ```sh
    bash install.sh
      ```
   P.S.: you might need to provide permissions with `chmod +x install.sh` command

   This will do the following:
   - Clone the repo (create a folder and download all needed files within it)
   - (Optional) Check if files have been downloaded properly (search for 'Uncomment' commentary in `install.sh`)
   - Add the executable to your resource `.rc` file if you are using Bash or Zsh shell. If you use a different shell you will have to add the executable manually  to your equivalent of resource `.rc` file (`.bashrc` in Bash, `.zshrc` in Zsh) of your shell (see step 2 in 'Manual Installation')

3. (Optional) You might need to give some permissions to the executable with:
   ```
   chmod +x /abs/path/to/PokeASCII/executable
   ```

4. Restart your terminal, or run:
   ```
   source YOUR_RESOURCE_FILE_HERE
   ```
   where `YOUR_RESOURCE_FILE_HERE` is `.zshrc` for Zsh and `.bashrc` for Bash shell.


#### II. Manual way
1. Clone the repo
    ```
    git clone github.com/p4nchit0z/PokeASCIILogin
    ```

2. Add the cloned executable to your resource `.rc` file (`.bashrc` for Bash shell, `.zshrc` for ZSH shell and so on...)

   ```
   echo "/abs/path/to/PokeASCII/executable" >> $HOME/"YOUR_RC_FILE_HERE"
   ```
   for example, if your user is `MyUsername` and you use `ZSH` shell:
   ```
   echo "/home/MyUsername/PokeASCIILogin/PokeASCIILogin" >> $HOME/.zshrc
   ```
   where `"/home/MyUsername/PokeASCIILogin/PokeASCIILogin"` is the absolute path to PokeASCIILogin executable.

3. (Optional) You might need to give some permissions to the executable with:
   ```
   chmod +x /abs/path/to/PokeASCII/executable
   ```

4. Restart your terminal, or run:
   ```
   source YOUR_RESOURCE_FILE_HERE
   ```
   where in my example this is:
   ```
   source $HOME/.zshrc
   ```
#### III. Build with Go (if installed)
1. Clone the repository
   ```
   git clone github.com/p4nchit0z/PokeASCIILogin
   ```
2. Build the executable
   ```
   cd PokeASCIILogin
   go build -o PokeASCIILogin
   ```
3. Add executable to your .rc file as explained in previous installations