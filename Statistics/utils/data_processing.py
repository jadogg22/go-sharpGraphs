import pandas as pd
import re

def extract_state(loc):
    m = re.search(r",\s*([A-Z]{2})$", str(loc))
    return m.group(1) if m else None

def load_and_prep_data(csv_file):
    df = pd.read_csv(csv_file)
    df = df[df["total_miles"] > 0]
    df = df[df["total_revenue"] > 0]
    df["bill_date"] = pd.to_datetime(df["bill_date"])
    df["OriginState"] = df["origin_state"]
    df["DestState"] = df["dest_state"]
    df = df[df["OriginState"] != df["DestState"]]
    df["RevenuePerMile"] = df["total_revenue"] / df["total_miles"]
    df["empty_pct"] = df["empty_miles"] / df["total_miles"]
    df["Route"] = df["origin_state"] + " → " + df["dest_state"]
    df["Brokered"] = df["customer_category"].apply(
        lambda x: "Brokered" if x == "BRK - BROKER" else "Direct"
    )
    return df

def analyze_lane_profitability(df):
    """
    Accurately analyzes lane profitability by pairing outbound and inbound legs
    without creating a Cartesian product.
    """
    # 1. Calculate stats for each one-way lane
    lane_stats = df.groupby(['OriginState', 'DestState']).agg(
        AvgRevPerMile=('RevenuePerMile', 'mean'),
        TripCount=('order_id', 'count'),
        AvgEmptyPct=('empty_pct', 'mean')
    ).reset_index()

    # 2. Create a "reverse" lane to merge on
    utah_outbound = lane_stats[lane_stats['OriginState'] == 'UT']
    utah_inbound = lane_stats[lane_stats['DestState'] == 'UT']

    # 3. Merge outbound from UT with inbound to UT
    merged_lanes = pd.merge(
        utah_outbound,
        utah_inbound,
        left_on='DestState',
        right_on='OriginState',
        suffixes=('_Outbound', '_Inbound')
    )

    # 4. Calculate combined metrics
    merged_lanes['AvgRoundTripRevenue'] = (merged_lanes['AvgRevPerMile_Outbound'] + merged_lanes['AvgRevPerMile_Inbound']) / 2
    merged_lanes['TotalTrips'] = merged_lanes['TripCount_Outbound'] + merged_lanes['TripCount_Inbound']
    merged_lanes['AvgEmptyPct'] = ((merged_lanes['AvgEmptyPct_Outbound'] * merged_lanes['TripCount_Outbound']) + \
                                  (merged_lanes['AvgEmptyPct_Inbound'] * merged_lanes['TripCount_Inbound'])) / \
                                 (merged_lanes['TripCount_Outbound'] + merged_lanes['TripCount_Inbound'])

    # 5. Clean up the dataframe
    final_lanes = merged_lanes[[
        'DestState_Outbound',
        'AvgRevPerMile_Outbound',
        'TripCount_Outbound',
        'AvgRevPerMile_Inbound',
        'TripCount_Inbound',
        'AvgRoundTripRevenue',
        'TotalTrips',
        'AvgEmptyPct'
    ]].rename(columns={'DestState_Outbound': 'DestinationState'})

    print(f"Lane analysis complete. Found {len(final_lanes)} round-trip lanes originating from UT.")
    return final_lanes
