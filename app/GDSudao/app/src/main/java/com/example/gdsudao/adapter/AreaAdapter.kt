package com.example.gdsudao.adapter

import android.content.Context
import android.content.Intent
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ActionMenuView
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.recyclerview.widget.RecyclerView
import com.example.gdsudao.DetalhesActivity
import com.example.gdsudao.activity.CadastroAreaActivity
import com.example.gdsudao.activity.MenuActivity
import com.example.gdsudao.model.Area
import kotlin.coroutines.coroutineContext

class AreaAdapter(context: Context, private val list: List<Area>) : RecyclerView.Adapter<AreaViewHolder>(){

    private val context: Context

    init {
        this.context = context
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): AreaViewHolder {
        val inflater = LayoutInflater.from(parent.context)
        return AreaViewHolder(inflater, parent)
    }

    override fun onBindViewHolder(holder: AreaViewHolder, position: Int) {
        val area : Area = list[position]
        holder.bind(area)

        // Funcionou
        holder.itemView.setOnClickListener{
            // ToDO: Abrir a tela de detalhes e passar o item
            var intent = Intent(context, DetalhesActivity::class.java)
            intent.putExtra("item", position)
            context.startActivity(intent)
        }
    }

    override fun getItemCount(): Int = list.size




}