package com.example.gdsudao.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.example.gdsudao.model.Area

class AreaAdapter(private val list: List<Area>) : RecyclerView.Adapter<AreaViewHolder>() {
    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): AreaViewHolder {
        val inflater = LayoutInflater.from(parent.context)
        return AreaViewHolder(inflater, parent)
    }

    override fun onBindViewHolder(holder: AreaViewHolder, position: Int) {
        val area : Area = list[position]
        holder.bind(area)
    }

    override fun getItemCount(): Int = list.size

}