import pandas as pd
from sklearn.preprocessing import StandardScaler

def calculate_custom_lane_score(lane_data):
    """
    Calculates a custom Lane Quality Score based on standardized features.
    """
    # Select features for the custom score
    features = ['AvgRevPerMile_Outbound', 'AvgRevPerMile_Inbound', 'TotalTrips', 'AvgEmptyPct']
    x = lane_data[features]

    # Standardize the features
    scaler = StandardScaler()
    x_scaled = scaler.fit_transform(x)

    # Create a DataFrame from scaled features for easier manipulation
    scaled_df = pd.DataFrame(x_scaled, columns=features, index=lane_data.index)

    # Invert the standardized 'AvgEmptyPct' so that a higher value is better
    scaled_df['AvgEmptyPct'] = scaled_df['AvgEmptyPct'] * -1

    # Define weights for each feature
    weights = {
        'AvgRevPerMile_Outbound': 1,
        'AvgRevPerMile_Inbound': 1,
        'TotalTrips': 1,
        'AvgEmptyPct': 0.5  # Reduced weight for empty percentage
    }

    # Calculate the custom score by summing the weighted standardized features
    lane_data['LaneQualityScore'] = (scaled_df['AvgRevPerMile_Outbound'] * weights['AvgRevPerMile_Outbound'] +
                                    scaled_df['AvgRevPerMile_Inbound'] * weights['AvgRevPerMile_Inbound'] +
                                    scaled_df['TotalTrips'] * weights['TotalTrips'] +
                                    scaled_df['AvgEmptyPct'] * weights['AvgEmptyPct'])

    print("Custom Lane Quality Score calculation complete.")

    # --- Print example calculations for verification ---
    print("\n--- Example Lane Quality Score Calculations ---")
    # Sort by score to pick top and bottom examples
    lane_data_sorted = lane_data.sort_values(by='LaneQualityScore', ascending=False)

    # Select a few examples (top 3 and bottom 3)
    examples = pd.concat([lane_data_sorted.head(3), lane_data_sorted.tail(3)])

    for index, row in examples.iterrows():
        print(f"\nLane: {row['DestinationState']}")
        print(f"  Original Outbound RPM: {row['AvgRevPerMile_Outbound']:.2f}")
        print(f"  Original Inbound RPM: {row['AvgRevPerMile_Inbound']:.2f}")
        print(f"  Original Total Trips: {row['TotalTrips']}")
        print(f"  Original Avg Empty Pct: {row['AvgEmptyPct']:.2f}%")

        # Get standardized values for this row
        scaled_outbound = scaled_df.loc[index, 'AvgRevPerMile_Outbound']
        scaled_inbound = scaled_df.loc[index, 'AvgRevPerMile_Inbound']
        scaled_trips = scaled_df.loc[index, 'TotalTrips']
        scaled_empty_pct = scaled_df.loc[index, 'AvgEmptyPct']

        print(f"  Standardized Outbound RPM: {scaled_outbound:.2f}")
        print(f"  Standardized Inbound RPM: {scaled_inbound:.2f}")
        print(f"  Standardized Total Trips: {scaled_trips:.2f}")
        print(f"  Standardized Avg Empty Pct (Inverted): {scaled_empty_pct:.2f}")
        
        # Calculate the weighted sum for the example
        calculated_score = (scaled_outbound * weights['AvgRevPerMile_Outbound'] +
                            scaled_inbound * weights['AvgRevPerMile_Inbound'] +
                            scaled_trips * weights['TotalTrips'] +
                            scaled_empty_pct * weights['AvgEmptyPct'])

        print(f"  Calculated Score (weighted sum): {calculated_score:.2f}")
        print(f"  Final Lane Quality Score: {row['LaneQualityScore']:.2f}")
    print("--------------------------------------------------")

    return lane_data

def analyze_small_customer_performance(df, lane_data):
    """
    Identifies small customers and analyzes the performance of the lanes they use.
    """
    # 1. Identify small customers (fewer than 200 loads)
    customer_trip_counts = df['Customer'].value_counts()
    small_customers = customer_trip_counts[customer_trip_counts < 200].index.tolist()

    # 2. Filter for trips by small customers
    small_customer_df = df[df['Customer'].isin(small_customers)]

    # 3. Merge with lane quality score data
    # We need to merge on the destination state
    customer_lane_performance = pd.merge(
        small_customer_df, 
        lane_data, 
        left_on='DestState', 
        right_on='DestinationState',
        how='left'
    )

    # 4. Select relevant columns and remove duplicates
    result = customer_lane_performance[['Customer', 'DestState', 'LaneQualityScore']].drop_duplicates()

    # 5. Sort by LaneQualityScore to find the worst-performing lanes
    result = result.sort_values(by='LaneQualityScore', ascending=True)

    # Save the report to a CSV
    result.to_csv("reports/small_customer_lane_performance.csv", index=False)
    print("Saved: reports/small_customer_lane_performance.csv")

    return result