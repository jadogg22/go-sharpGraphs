
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import os
from matplotlib.colors import LogNorm

def make_state_bar_chart(df_subset, title, filename, target_revenue=2.79):
    outbound = df_subset.groupby("OriginState")["RevenuePerMile"].mean().rename("OutboundRevPerMile")
    inbound = df_subset.groupby("DestState")["RevenuePerMile"].mean().rename("InboundRevPerMile")
    state_revenue = pd.concat([outbound, inbound], axis=1).fillna(0)

    state_revenue = state_revenue.sort_values("OutboundRevPerMile", ascending=False)
    states = state_revenue.index
    x = range(len(states))
    bar_width = 0.4

    plt.figure(figsize=(14, 6))
    plt.bar(x, state_revenue["OutboundRevPerMile"], width=bar_width, label="Outbound", align='center')
    plt.bar([i + bar_width for i in x], state_revenue["InboundRevPerMile"], width=bar_width, label="Inbound", align='center')

    # Add target revenue line
    plt.axhline(y=target_revenue, color='red', linestyle='--', linewidth=1, label=f'Target: ${target_revenue:.2f}')

    plt.xlabel("State")
    plt.ylabel("Avg Revenue per Mile")
    plt.title(title)
    plt.xticks([i + bar_width/2 for i in x], states, rotation=45)
    plt.legend()
    plt.tight_layout()
    plt.savefig(f"reports/{filename}")
    plt.close()

    # Also save to CSV
    state_revenue.reset_index().to_csv(f"reports/{filename.replace('.png', '.csv')}", index=False)

def plot_lane_frequency(df):
    # Filter out brokered loads
    df_direct = df[df["Brokered"] == "Direct"]

    lane_stats = df_direct.groupby("Route").agg(
        Frequency=("order_id", "count"),
        AvgRevPerMile=("RevenuePerMile", "mean")
    ).reset_index()

    lane_stats.to_csv("reports/lane_summary.csv", index=False)

    plt.figure(figsize=(12, 6))
    sns.scatterplot(
        data=lane_stats,
        x="Frequency",
        y="AvgRevPerMile",
        size="Frequency",
        sizes=(20, 100),
        alpha=0.8
    )
    plt.title("Lane Frequency vs Revenue per Mile (Direct Loads)")
    plt.xlabel("Number of Loads")
    plt.ylabel("Avg Revenue per Mile")
    plt.tight_layout()
    plt.savefig("reports/lane_scatter.png")
    plt.close()

def plot_brokered_comparison(df):
    brokered_stats = df.groupby("Brokered").agg(
        AvgRevPerMile=("RevenuePerMile", "mean"),
        EmptyPct=("empty_pct", "mean"),
        TotalOrders=("order_id", "count"),
        TotalRevenue=("total_revenue", "sum")
    ).reset_index()

    brokered_stats.to_csv("reports/brokered_summary.csv", index=False)

    fig, axs = plt.subplots(1, 2, figsize=(14, 5))

    melted = brokered_stats.melt(id_vars="Brokered", value_vars=["AvgRevPerMile", "EmptyPct"])
    sns.barplot(data=melted, x="variable", y="value", hue="Brokered", ax=axs[0])
    axs[0].set_title("Rev per Mile & Empty % by Brokered Status")
    axs[0].set_ylabel("")
    axs[0].set_xlabel("")

    sns.barplot(data=brokered_stats, x="Brokered", y="TotalOrders", ax=axs[1])
    axs[1].set_title("Total Orders by Brokered Status")
    axs[1].set_ylabel("")

    plt.tight_layout()
    plt.savefig("reports/brokered_comparison.png")
    plt.close()

def plot_lane_profitability_scatter(lane_data, target_revenue=2.79):
    """
    Plots the outbound vs. inbound revenue per mile for each lane.
    """
    lane_data.to_csv("reports/lane_profitability_scatter.csv", index=False)
    plt.figure(figsize=(12, 8))
    sns.scatterplot(
        data=lane_data,
        x="AvgRevPerMile_Outbound",
        y="AvgRevPerMile_Inbound",
        size="TotalTrips",
        hue="DestinationState",
        palette="tab20",
        legend=False,
        edgecolor='black',
        alpha=0.7,
        sizes=(40, 400)
    )

    # Add target lines
    plt.axhline(y=target_revenue, color='red', linestyle='--', linewidth=1, label=f'Target: ${target_revenue:.2f}')
    plt.axvline(x=target_revenue, color='red', linestyle='--', linewidth=1)

    for _, row in lane_data.iterrows():
        plt.text(
            row["AvgRevPerMile_Outbound"],
            row["AvgRevPerMile_Inbound"] + 0.02,
            row["DestinationState"],
            fontsize=9,
            ha="center"
        )

    plt.title("Outbound vs. Inbound Revenue per Mile by Lane")
    plt.xlabel("Average Outbound Revenue per Mile")
    plt.ylabel("Average Inbound Revenue per Mile")
    plt.grid(True)
    plt.legend()
    plt.tight_layout()

    plot_path = "reports/lane_profitability_scatter.png"
    plt.savefig(plot_path, dpi=300)
    plt.close()
    print(f"Saved lane profitability scatter plot to: {plot_path}")

def plot_avg_lane_revenue(lane_data, target_revenue=2.79):
    """
    Plots the average round-trip revenue per mile for each lane.
    """
    lane_data = lane_data.sort_values(by='AvgRoundTripRevenue', ascending=False)
    lane_data.to_csv("reports/avg_lane_revenue.csv", index=False)

    # Create a logarithmic color map
    norm = LogNorm(lane_data['TotalTrips'].min(), lane_data['TotalTrips'].max())
    sm = plt.cm.ScalarMappable(cmap="viridis", norm=norm)
    sm.set_array([])

    # Plot
    fig, ax = plt.subplots(figsize=(15, 8))
    bars = ax.bar(lane_data['DestinationState'], lane_data['AvgRoundTripRevenue'], color=sm.to_rgba(lane_data['TotalTrips']))

    # Add labels to the bars
    for bar in bars:
        height = bar.get_height()
        ax.annotate(f'{lane_data.loc[lane_data["AvgRoundTripRevenue"] == height, "TotalTrips"].iloc[0]}',
                    xy=(bar.get_x() + bar.get_width() / 2, height),
                    xytext=(0, 3),  # 3 points vertical offset
                    textcoords="offset points",
                    ha='center', va='bottom')

    # Add color bar
    cbar = fig.colorbar(sm, ax=ax)
    cbar.set_label('Total Round Trips (Log Scale)')

    # Horizontal target line
    ax.axhline(y=target_revenue, color='red', linestyle='--', linewidth=2, label=f'Target: ${target_revenue:.2f}')

    # Chart styling
    ax.set_ylabel('Average Round-Trip $/Mile')
    ax.set_xlabel('Destination State')
    ax.set_title('Average Round-Trip Revenue per Mile by Lane')
    plt.xticks(rotation=45)
    ax.legend()
    plt.tight_layout()

    # Save chart
    plt.savefig('reports/avg_lane_revenue.png')
    plt.close()
    print("Saved: reports/avg_lane_revenue.png")

def plot_revenue_by_day(df):
    """
    Plots the average revenue per mile over time.
    """
    daily_revenue = df.groupby('bill_date')['RevenuePerMile'].mean().reset_index()
    daily_revenue.to_csv("reports/revenue_by_day.csv", index=False)

    plt.figure(figsize=(15, 7))
    plt.plot(daily_revenue['bill_date'], daily_revenue['RevenuePerMile'], marker='o', linestyle='-')

    # Chart styling
    plt.title('Average Revenue per Mile Over Time')
    plt.ylabel('Average Revenue per Mile')
    plt.xlabel('Bill Date')
    plt.grid(True)
    plt.tight_layout()

    # Save chart
    plt.savefig('reports/revenue_by_day.png')
    plt.close()
    print("Saved: reports/revenue_by_day.png")

def plot_custom_lane_scores(lane_data):
    """
    Plots the custom Lane Quality Score for each lane.
    """
    lane_data_sorted = lane_data.sort_values(by='LaneQualityScore', ascending=False)
    lane_data_sorted.to_csv("reports/custom_lane_quality_scores.csv", index=False)

    plt.figure(figsize=(15, 8))
    sns.barplot(x='DestinationState', y='LaneQualityScore', data=lane_data_sorted, palette='viridis')

    plt.title('Custom Lane Quality Score by Destination State (Higher is Better)')
    plt.xlabel('Destination State')
    plt.ylabel('Custom Lane Quality Score')
    plt.xticks(rotation=45)
    plt.tight_layout()

    plt.savefig('reports/custom_lane_quality_scores.png')
    plt.close()
    print("Saved: reports/custom_lane_quality_scores.png")
