Welcome to the Epidemic and Herd Immunity Simulator!
The system is designed to be fairly easy to use, and is hopefully fun
to play around with.

INSTALLATION INSTRUCTIONS:
Please install the entire contents of the .tar file into your go/src/ directory.

IMPORTANT -- If you do not want the program to be run as "dis.exe", 
please delete "dis.exe" from the installation folder, rename 
the /dis folder to whatever you want, and re-run the "go build" command
from the command line.




RUNNING INSTRUCTIONS:
NOTE: FOR ALL USER INPUTS, PLEASE DO NOT INCLUDE WHITESPACES

To run the program, simply run 'dis.exe', or navigate to your command line
and run the program according to your OS. Please do not delete or alter
the /pathogens or /progression directories.

1) You will first be asked to load a .PATHOGEN file to run the simulations.
These are stored in the /pathogens directory, which contains a number
of sample pathogens for your convenience. When entering this, be sure
to include the file extension in your input ("measles.PATHOGEN" not
"measles"). 

If you want to experiment with creating your own disease, enter "CUSTOM"
where prompted and follow the prompts. Your new pathogen should be saved
into the /pathogens directory.

2) You will then be prompted to enter a population size for your community,
this should be an integer number greater than 0. The program should run
extremely quickly with population sizes up to 150,000. Population values
>1,000,000 will run slowly and will also be impossible to visualize later.

3) You will then be prompted to enter a vaccination rate as a positive
integer between 0 and 100

4) You will finally be prompted to enter a number of individuals to start
the infection. I recommend starting with a number in between 1 and 5, 
depending on the size of your population.

The simulation will then run, and write a .gif animation showing the
progression of your disease. It will also write a .txt file containing
the statistics from your epidemic.

The /progression directory will store all the images contained in the
animation in order. That is, the state of the infection at each timestep.

NOTE: The /progression directory will overwrite every time the simulation
is run, so be sure to move or rename your files if you do not want to lose
them.

I had a great time designing and implementing this project, and I hope
you enjoy messing around with it!


