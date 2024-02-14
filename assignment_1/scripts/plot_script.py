import matplotlib.pyplot as plt

with open('avg_response.txt', 'r') as file:
    lines = file.readlines()

x, y = zip(*(line.split() for line in lines))

x = [int(i) for i in x]
y = [float(i) for i in y]

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(12, 6))
ax1.plot(x, y, marker='o')
ax1.set_title('Avg Response Time')
ax1.set_xlabel('Number of containers')
ax1.set_ylabel('Response Time (seconds)')

table_data = []
for i, (xi, yi) in enumerate(zip(x, y)):
    table_data.append([xi, yi])

table = ax2.table(cellText=table_data, loc='center', colLabels=["Number of containers", "Response Time (seconds)"], fontsize=8)
table.auto_set_font_size(False)
ax2.axis('off') 

plt.savefig('avg_response_plot.png', bbox_inches='tight')
fig.savefig('avg_response_table.png', bbox_inches='tight')

plt.close()