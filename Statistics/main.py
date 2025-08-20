import os
import sys
import seaborn as sns
from utils.data_processing import load_and_prep_data, analyze_lane_profitability
from utils.plotting import (
    make_state_bar_chart,
    plot_lane_frequency,
    plot_brokered_comparison,
    plot_lane_profitability_scatter,
    plot_avg_lane_revenue,
    plot_revenue_by_day,
    plot_custom_lane_scores
)
from utils.advanced_analysis import calculate_custom_lane_score
from utils.reporting import generate_pdf_report

def main(csv_file, pdf_output_path):
    # --- Setup ---
    sns.set(style="whitegrid")
    os.makedirs("reports", exist_ok=True)

    # --- Load and Prep Data ---
    df = load_and_prep_data(csv_file)

    # === 1. General Analyses ===
    make_state_bar_chart(df, "All Loads: Avg Revenue per Mile by State", "state_revenue_all.png")
    df_brokered = df[df["Brokered"] == "Brokered"]
    make_state_bar_chart(df_brokered, "Brokered Only: Avg Revenue per Mile by State", "state_revenue_brokered.png")

    # === 2. Lane Profitability Analysis ===
    lane_data = analyze_lane_profitability(df)
    print("Generated lane profitability summary.")

    plot_lane_profitability_scatter(lane_data)
    plot_avg_lane_revenue(lane_data)

    # === 3. Custom Lane Quality Score ===
    lane_data_with_custom_score = calculate_custom_lane_score(lane_data.copy()) # Use a copy to avoid modifying original lane_data
    plot_custom_lane_scores(lane_data_with_custom_score)

    # === 4. Generate PDF Report ===
    report_content = [
        {
            "title": "Custom Lane Quality Score",
            "text": """
This section introduces the Custom Lane Quality Score, a metric designed to provide a holistic view of lane performance. 

How it's calculated:
1.  Key Metrics: We analyze Average Outbound Revenue per Mile, Average Inbound Revenue per Mile, Total Trips, and Average Empty Miles Percentage.
2.  Standardization: Each metric is standardized to ensure fair comparison, regardless of its original scale.
3.  Weighted Calculation: The final score is a weighted sum of these standardized metrics. Revenue and trip volume are given full weight, while Empty Mile Percentage is included with a lower weight to act as a penalty without dominating the score.

A higher score indicates a more desirable lane, balancing revenue, volume, and operational efficiency.""",
            "image": "custom_lane_quality_scores.png"
        },
        {
            "title": "Average Lane Revenue",
            "text": """This chart provides a clear comparison of the average round-trip revenue per mile for each destination state.

How to read this chart:
- The height of each bar represents the average revenue per mile for that lane. Higher bars are better.
- The number on top of each bar is the total number of round trips for that lane. This provides context for the revenue average. A high-revenue lane with very few trips may be an outlier, while one with many trips is a reliable source of income.
- The color of the bar also indicates the total number of trips, with darker shades representing higher volume. This allows for a quick visual assessment of which lanes are both profitable and high-volume.""",
            "image": "avg_lane_revenue.png"
        },
        {
            "title": "State-Level Revenue Analysis",
            "text": "This section provides a high-level overview of average revenue per mile, broken down by state for both all loads and brokered-only loads.",
            "image": "state_revenue_all.png"
        },
        {
            "title": "Lane Profitability",
            "text": """This chart provides a detailed look at the profitability of each lane, considering both inbound and outbound revenue. The red lines represent the average revenue per mile across all lanes.

Your goal is to have lanes in the TOP-RIGHT quadrant, which indicates that both the outbound and inbound legs of the trip are performing above average. Lanes in the bottom-left quadrant are underperforming in both directions and should be investigated.""",
            "image": "lane_profitability_scatter.png"
        }
    ]

    generate_pdf_report(report_content, output_path=pdf_output_path)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python main.py <input_csv_path> <output_pdf_path>")
        sys.exit(1)
    
    input_csv = sys.argv[1]
    output_pdf = sys.argv[2]
    main(input_csv, output_pdf)
