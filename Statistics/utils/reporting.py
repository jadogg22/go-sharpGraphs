import os
import matplotlib.pyplot as plt
from matplotlib.backends.backend_pdf import PdfPages

def generate_pdf_report(report_content, output_path):
    with PdfPages(output_path) as pdf:
        for item in report_content:
            # Create a new page for each item
            fig = plt.figure(figsize=(8.5, 11))

            # Add title and text
            plt.figtext(0.5, 0.95, item['title'], ha='center', va='top', fontsize=16, weight='bold')
            plt.figtext(0.1, 0.85, item['text'], ha='left', va='top', wrap=True, fontsize=12)

            if 'image' in item:
                # Add image
                img_path = os.path.join("reports", item['image'])
                if os.path.exists(img_path):
                    img_data = plt.imread(img_path)
                    img_ax = fig.add_axes([0.1, 0.1, 0.8, 0.7]) # Adjust position and size as needed
                    img_ax.imshow(img_data)
                    img_ax.axis('off')
                else:
                    print(f"Warning: Image file not found at {img_path}")
            elif 'table' in item:
                # Add table
                table_data = item['table']
                table = plt.table(cellText=table_data.values, colLabels=table_data.columns, loc='center', cellLoc='center')
                table.auto_set_font_size(False)
                table.set_fontsize(10)
                table.scale(1, 1.5)
                plt.axis('off')

            pdf.savefig(fig, dpi=300)
            plt.close()