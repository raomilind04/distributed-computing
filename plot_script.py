
import matplotlib.pyplot as plt

# Reading data from the file
with open('data.txt', 'r') as file:
    lines = file.readlines()

# Splitting data into x and y values
x, y = zip(*(line.split() for line in lines))

# Converting string data to integers
x = [int(i) for i in x]
y = [int(i) for i in y]

# Creating the plot
plt.figure(figsize=(10, 6))
plt.plot(x, y, marker='o')

# Adding title and labels
plt.title('Total Response Times')
plt.xlabel('X axis')
plt.ylabel('Response Time (seconds)')

# Saving the plot
plt.savefig('plot.png')
